import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MultiClusterRelationCreateComponent } from './multi-cluster-relation-create.component';

describe('MultiClusterRelationCreateComponent', () => {
  let component: MultiClusterRelationCreateComponent;
  let fixture: ComponentFixture<MultiClusterRelationCreateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MultiClusterRelationCreateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MultiClusterRelationCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
