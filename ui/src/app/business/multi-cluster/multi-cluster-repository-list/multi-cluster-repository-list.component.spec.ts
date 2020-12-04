import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MultiClusterRepositoryListComponent } from './multi-cluster-repository-list.component';

describe('MultiClusterRepositoryListComponent', () => {
  let component: MultiClusterRepositoryListComponent;
  let fixture: ComponentFixture<MultiClusterRepositoryListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MultiClusterRepositoryListComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MultiClusterRepositoryListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
