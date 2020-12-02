import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MultiClusterRepositoryDetailComponent } from './multi-cluster-repository-detail.component';

describe('MultiClusterRepositoryDetailComponent', () => {
  let component: MultiClusterRepositoryDetailComponent;
  let fixture: ComponentFixture<MultiClusterRepositoryDetailComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MultiClusterRepositoryDetailComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MultiClusterRepositoryDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
